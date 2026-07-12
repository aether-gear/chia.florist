import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Loader2, Save, Trash2, UploadCloud, AlertTriangle, CheckCircle2, Image as ImageIcon } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { Badge } from '../../components/ui/badge';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from '../../components/ui/dialog';
import { useEditProductViewModel, type EditProductForm } from '../../viewmodels/useEditProductViewModel';

export default function EditProductPage() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();

  const {
    product,
    loading,
    error,
    saving,
    saveError,
    uploading,
    uploadError,
    deleting,
    deleteError,
    updateProduct,
    uploadImage,
    deleteProduct,
  } = useEditProductViewModel(slug);

  const [formData, setFormData] = useState<EditProductForm>({
    sku: '',
    name: '',
    description: '',
    is_available: true,
    status: 'active',
    price: 0,
    weight: null,
  });

  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [uploadSuccess, setUploadSuccess] = useState(false);

  // Sync state when product loads
  useEffect(() => {
    if (product) {
      setFormData({
        sku: product.sku || '',
        name: product.name || '',
        description: product.description || '',
        is_available: product.is_available,
        status: product.status || 'active',
        price: product.price || 0,
        weight: product.weight !== undefined && product.weight !== null ? Number(product.weight) : null,
      });
    }
  }, [product]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value, type } = e.target as HTMLInputElement;
    const isCheckbox = type === 'checkbox';

    setFormData((prev) => ({
      ...prev,
      [name]: isCheckbox ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaveSuccess(false);
    const success = await updateProduct(formData);
    if (success) {
      setSaveSuccess(true);
      window.scrollTo({ top: 0, behavior: 'smooth' });
      setTimeout(() => setSaveSuccess(false), 3000);
    }
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      setUploadSuccess(false);
      const success = await uploadImage(file);
      if (success) {
        setUploadSuccess(true);
        setTimeout(() => setUploadSuccess(false), 3000);
      }
    }
  };

  const handleDeleteConfirm = async () => {
    await deleteProduct();
  };

  if (loading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="flex h-[50vh] flex-col items-center justify-center space-y-4">
        <p className="text-destructive font-medium">{error || 'Product not found'}</p>
        <Button variant="outline" onClick={() => navigate('/products')}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Products
        </Button>
      </div>
    );
  }

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-4 p-8 pt-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <Button variant="outline" size="icon" onClick={() => navigate('/products')}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <div className="flex items-center gap-2">
                <h2 className="text-3xl font-bold tracking-tight">Edit Product</h2>
                <Badge variant={product.status === 'active' ? 'default' : 'secondary'}>
                  {product.status}
                </Badge>
              </div>
              <p className="text-muted-foreground">Manage catalog details, images, and inventory for {product.name}</p>
            </div>
          </div>
        </div>

        {/* Alerts */}
        {saveSuccess && (
          <div className="flex items-center gap-2 p-4 text-sm text-emerald-600 bg-emerald-50 rounded-lg border border-emerald-100 shadow-sm animate-in fade-in slide-in-from-top-1">
            <CheckCircle2 className="h-5 w-5 shrink-0" />
            <span>Product details successfully saved and updated.</span>
          </div>
        )}
        {saveError && (
          <div className="flex items-center gap-2 p-4 text-sm text-red-600 bg-red-50 rounded-lg border border-red-100 shadow-sm">
            <AlertTriangle className="h-5 w-5 shrink-0" />
            <span>{saveError}</span>
          </div>
        )}
        {uploadSuccess && (
          <div className="flex items-center gap-2 p-4 text-sm text-emerald-600 bg-emerald-50 rounded-lg border border-emerald-100 shadow-sm">
            <CheckCircle2 className="h-5 w-5 shrink-0" />
            <span>Product image uploaded successfully.</span>
          </div>
        )}
        {uploadError && (
          <div className="flex items-center gap-2 p-4 text-sm text-red-600 bg-red-50 rounded-lg border border-red-100 shadow-sm">
            <AlertTriangle className="h-5 w-5 shrink-0" />
            <span>{uploadError}</span>
          </div>
        )}
        {deleteError && (
          <div className="flex items-center gap-2 p-4 text-sm text-red-600 bg-red-50 rounded-lg border border-red-100 shadow-sm">
            <AlertTriangle className="h-5 w-5 shrink-0" />
            <span>{deleteError}</span>
          </div>
        )}

        <div className="grid gap-6 md:grid-cols-3">
          {/* Form Column */}
          <div className="md:col-span-2">
            <Card className="border border-slate-100 shadow-sm">
              <CardHeader>
                <CardTitle>Product Details</CardTitle>
                <CardDescription>Update general fields, pricing, and availability</CardDescription>
              </CardHeader>
              <CardContent>
                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="sku">SKU</Label>
                      <Input
                        id="sku"
                        name="sku"
                        value={formData.sku}
                        onChange={handleChange}
                        required
                        className="bg-transparent"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="name">Product Name</Label>
                      <Input
                        id="name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                        className="bg-transparent"
                      />
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="description">Description (Optional)</Label>
                    <textarea
                      id="description"
                      name="description"
                      className="flex min-h-[120px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                      value={formData.description || ''}
                      onChange={handleChange}
                      placeholder="Enter a descriptive detail for the product..."
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="price">Price (IDR)</Label>
                      <Input
                        id="price"
                        name="price"
                        type="number"
                        min="0"
                        value={formData.price}
                        onChange={handleChange}
                        required
                        className="bg-transparent"
                      />
                    </div>
                    <div className="space-y-2">
                      <Label htmlFor="weight">Weight (kg) (Optional)</Label>
                      <Input
                        id="weight"
                        name="weight"
                        type="number"
                        step="0.01"
                        min="0"
                        value={formData.weight || ''}
                        onChange={handleChange}
                        className="bg-transparent"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label htmlFor="status">Status</Label>
                      <select
                        id="status"
                        name="status"
                        className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                        value={formData.status}
                        onChange={handleChange}
                      >
                        <option value="active">Active</option>
                        <option value="inactive">Inactive</option>
                        <option value="archived">Archived</option>
                      </select>
                    </div>

                    <div className="space-y-2 flex flex-col justify-end">
                      <div className="flex items-center space-x-2 h-9">
                        <input
                          type="checkbox"
                          id="is_available"
                          name="is_available"
                          className="h-4 w-4 rounded border-gray-300 accent-indigo-600"
                          checked={formData.is_available}
                          onChange={handleChange}
                        />
                        <Label htmlFor="is_available" className="font-normal cursor-pointer text-slate-700">
                          Available for sale
                        </Label>
                      </div>
                    </div>
                  </div>

                  <div className="flex justify-end border-t pt-4 gap-2">
                    <Button type="button" variant="outline" onClick={() => navigate('/products')} disabled={saving}>
                      Cancel
                    </Button>
                    <Button type="submit" disabled={saving}>
                      {saving ? (
                        <>
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving...
                        </>
                      ) : (
                        <>
                          <Save className="mr-2 h-4 w-4" /> Save Details
                        </>
                      )}
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </div>

          {/* Right Column: Images & Actions */}
          <div className="space-y-6">
            {/* Image Gallery & Upload */}
            <Card className="border border-slate-100 shadow-sm">
              <CardHeader>
                <CardTitle>Product Images</CardTitle>
                <CardDescription>Upload or view product images</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Current Banner / Main Image */}
                <div className="relative aspect-square w-full overflow-hidden rounded-lg border bg-slate-50 flex items-center justify-center">
                  {product.banner?.thumbnail ? (
                    <img
                      src={product.banner.thumbnail}
                      alt={product.name}
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <div className="flex flex-col items-center justify-center text-slate-400 p-4 text-center">
                      <ImageIcon className="h-10 w-10 mb-2 stroke-[1.5]" />
                      <span className="text-xs">No main image uploaded</span>
                    </div>
                  )}
                </div>

                {/* Image Gallery Grid */}
                {product.gallery && product.gallery.length > 0 && (
                  <div className="space-y-2">
                    <Label className="text-xs font-semibold text-slate-500 uppercase">Gallery</Label>
                    <div className="grid grid-cols-4 gap-2">
                      {product.gallery.map((img, idx) => (
                        <div key={idx} className="relative aspect-square rounded-md border overflow-hidden bg-slate-50">
                          <img
                            src={img.thumbnail || img.preview || img.detail || ''}
                            alt={`gallery-${idx}`}
                            className="h-full w-full object-cover"
                          />
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Upload Button */}
                <div className="pt-2">
                  <input
                    type="file"
                    id="product-image-upload"
                    accept="image/*"
                    onChange={handleFileChange}
                    className="hidden"
                    disabled={uploading}
                  />
                  <Label
                    htmlFor="product-image-upload"
                    className="flex flex-col items-center justify-center w-full h-28 border-2 border-dashed border-slate-200 rounded-lg cursor-pointer hover:bg-slate-50 hover:border-indigo-400 transition-all text-slate-500"
                  >
                    {uploading ? (
                      <div className="flex flex-col items-center">
                        <Loader2 className="h-6 w-6 animate-spin text-indigo-600 mb-1" />
                        <span className="text-xs font-medium">Uploading image...</span>
                      </div>
                    ) : (
                      <div className="flex flex-col items-center px-4 py-2 text-center">
                        <UploadCloud className="h-6 w-6 text-slate-400 mb-1" />
                        <span className="text-xs font-semibold">Upload Image</span>
                        <span className="text-[10px] text-slate-400 mt-0.5">Supports PNG, JPG, JPEG</span>
                      </div>
                    )}
                  </Label>
                </div>
              </CardContent>
            </Card>

            {/* Danger Zone */}
            <Card className="border border-red-100 bg-red-50/20 shadow-sm">
              <CardHeader>
                <CardTitle className="text-red-800 flex items-center gap-1.5 text-lg">
                  <AlertTriangle className="h-5 w-5 text-red-600 shrink-0" />
                  Danger Zone
                </CardTitle>
                <CardDescription className="text-red-700/80">
                  Permanently delete this product from the catalog. This action cannot be undone.
                </CardDescription>
              </CardHeader>
              <CardContent className="pt-0">
                <Button
                  type="button"
                  variant="destructive"
                  className="w-full bg-red-600 hover:bg-red-700 font-medium"
                  onClick={() => setIsDeleteDialogOpen(true)}
                  disabled={deleting}
                >
                  {deleting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Deleting...
                    </>
                  ) : (
                    <>
                      <Trash2 className="mr-2 h-4 w-4" /> Delete Product
                    </>
                  )}
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-red-600">
              <AlertTriangle className="h-5 w-5" />
              Confirm Deletion
            </DialogTitle>
            <DialogDescription className="py-2">
              Are you sure you want to permanently delete <strong>{product.name}</strong>? All association with inventory and shops will be removed. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex gap-2">
            <DialogClose asChild>
              <Button type="button" variant="outline">
                Cancel
              </Button>
            </DialogClose>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDeleteConfirm}
              disabled={deleting}
            >
              {deleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Yes, Delete Product
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
